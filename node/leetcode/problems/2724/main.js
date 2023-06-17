/**
 * @param {Array} arr
 * @param {Function} fn
 * @return {Array}
 */
const sortBy = function (arr, fn) {
  return arr.sort((a, b) => {
    return fn(a) - fn(b);
  });
};

// [1,2,3,4,5]
console.log(sortBy([5, 4, 1, 2, 3], (x) => x));
// [{"x": -1}, {"x": 0}, {"x": 1}]
console.log(sortBy([{ "x": 1 }, { "x": 0 }, { "x": -1 }], (x) => x.x));
// [[10, 1], [5, 2], [3, 4]]
console.log(sortBy([[3, 4], [5, 2], [10, 1]], (x) => x[1]));
