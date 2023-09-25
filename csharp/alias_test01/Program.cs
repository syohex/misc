public class Foo {
    public static void Main(string[] args) {
        IntGenerics a = new(10);
        Console.WriteLine($"a={a.Field}");
    }
}
