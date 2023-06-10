/**
 * @param {Promise} promise1
 * @param {Promise} promise2
 * @return {Promise}
 */
const addTwoPromises = async function(promise1, promise2) {
    const n1 = await promise1;
    const n2 = await promise2;

    return n1 + n2;
};

async function main() {
    {
        const promise1 = new Promise(resolve => setTimeout(() => resolve(2), 20));
        const promise2 = new Promise(resolve => setTimeout(() => resolve(5), 60));
        const ret = await addTwoPromises(promise1, promise2);
        // 7
        console.log(ret);
    }
    {
        const promise1 = new Promise(resolve => setTimeout(() => resolve(10), 20));
        const promise2 = new Promise(resolve => setTimeout(() => resolve(-12), 60));
        const ret = await addTwoPromises(promise1, promise2);
        // -2
        console.log(ret);
    }
}

main().catch(e => console.log(e));
