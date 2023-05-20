const ArrayWrapper = function (nums) {
  this.nums = nums;
};

ArrayWrapper.prototype.valueOf = function () {
  let ret = 0;
  for (const v of this.nums) {
    ret += v;
  }

  return ret;
};

ArrayWrapper.prototype.toString = function () {
  const s = this.nums.map((n) => n.toString()).join(",");
  return `[${s}]`;
};

{
  const obj1 = new ArrayWrapper([1, 2]);
  const obj2 = new ArrayWrapper([3, 4]);
  console.log(obj1 + obj2); // 10
}
{
  const obj1 = new ArrayWrapper([23, 98, 42, 70]);
  console.log(String(obj1)); // "[23,98,42,70]"
}

{
  const obj1 = new ArrayWrapper([]);
  const obj2 = new ArrayWrapper([]);
  console.log(obj1 + obj2); // 0
}
