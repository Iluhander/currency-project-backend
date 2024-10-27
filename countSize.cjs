const path = require('path');
const fs = require('fs');

const format = '   lines  cRate  name';

function formattedLog({ linesCount, cRate, fullPath }) {
  let rowLog = Array(format.length).fill(' ');
  rowLog[0] = '+';

  const data = [String(linesCount), String(cRate), String(fullPath)];
  const colsPositions = [3, 9, 16];

  for (let k = 0; k < data.length; ++k) {
    let curResPos = colsPositions[k];

    for (let i = 0; i < data[k].length; ++i) {
      rowLog[curResPos] = data[k][i];
      curResPos += 1;
    }
  }

  console.log(rowLog.join(''));
}

function walk(dir, exclusions) {
  let totalSize = 0;
  let totalCount = 0;
  let commentsCount = 0;

  fs.readdirSync(dir).forEach((file) => {
    const fullPath = path.join(dir, file);

    if (fs.lstatSync(fullPath).isDirectory()) {
      const subData = walk(fullPath, exclusions);
      
      totalSize += subData.totalSize;
      totalCount += subData.totalCount;
      commentsCount += subData.commentsCount;
    } else {
      if (exclusions.find((exclusion) => exclusion.test(fullPath))) {
        return;
      }

      const content = fs.readFileSync(fullPath, 'utf8');
      const splittedContent = content.split('\n');
      const linesCount = splittedContent.length;

      const commentLines = splittedContent.reduce((curCount, line) => {
        const trimmed = line.trim();

        if (
          trimmed.startsWith('*') ||
          trimmed.startsWith('/*') ||
          trimmed.indexOf('//') > -1
        ) {
          return curCount + 1;
        }

        return curCount;
      }, 0);

      // Updating the statistics.
      totalSize += linesCount;
      totalCount += 1;
      commentsCount += commentLines;

      const cRate = Math.round((commentLines / linesCount * 100));
      formattedLog({ linesCount, cRate, fullPath });
    }
  });

  return { totalSize, totalCount, commentsCount };
}

const exclusions = [/DS_Store/, /svg$/, /git/];

console.log(format);
const { totalSize, totalCount, commentsCount } = walk('./internal', exclusions);


console.log(`\n> ${totalSize} lines in ${totalCount} files`);
console.log(`> average: ${Math.round((totalSize / totalCount) * 100) / 100} per file`);
console.log(`> total cRate ${Math.round((commentsCount / totalSize * 10000)) / 100}%`);

