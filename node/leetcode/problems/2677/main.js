/**
 * @param {Array} arr
 * @param {number} size
 * @return {Array[]}
 */
const chunk = function (arr, size) {
  const ret = [];
  for (let i = 0; i < arr.length;) {
    const elem = [];
    for (let j = 0; j < size && i < arr.length; ++j) {
      elem.push(arr[i++]);
    }
    ret.push(elem);
  }
  return ret;
};

// [[1],[2],[3],[4],[5]]
console.log(chunk([1, 2, 3, 4, 5], 1));

// [[1,9,6], [3,2]]
console.log(chunk([1, 9, 6, 3, 2], 3));

// [[8,5,3,2,6]]
console.log(chunk([8, 5, 3, 2, 6], 6));

// []
console.log(chunk([], 1));
