var fibGenerator = function*() {
  yield 0;
  yield 1;

  let prev = 0;
  let cur = 1;

  while (true) {
    yield cur + prev;

    const tmp = prev;
    prev = cur;
    cur = cur + tmp;
  }
};

const gen = fibGenerator();
for (let i = 0; i < 10; ++i) {
  console.log(gen.next().value);
}
