const expect = function(val) {
    return {
        toBe: (v) => {
            if (v === val) {
                return true;
            }

            throw new Error("Not Equal");
        },
        notToBe: (v) => {
            if (v !== val) {
                return true;
            }

            throw new Error("Equal");
        }
    };
};

console.log(expect(5).toBe(5));
console.log(expect(5).notToBe(null));
