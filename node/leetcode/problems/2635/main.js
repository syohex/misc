var map = function (arr, fn) {
  let ret = [];
  for (let i = 0; i < arr.length; ++i) {
    ret.push(fn(arr[i], i));
  }
  return ret;
};

// [2, 3, 4]
console.log(map([1, 2, 3], (n) => n + 1));

// [1, 3, 5]
console.log(map([1, 2, 3], (n, i) => n + i));

// [42, 42, 42]
console.log(map([10, 20, 30], (_) => 42));
