const createCounter = function (init) {
  let val = init;
  return {
    increment: () => {
      val++;
      return val;
    },
    decrement: () => {
      val--;
      return val;
    },
    reset: () => {
      val = init;
      return val;
    },
  };
};

const counter = createCounter(5);
console.log(counter.increment()); // 6
console.log(counter.reset()); // 5
console.log(counter.decrement()); // 4
