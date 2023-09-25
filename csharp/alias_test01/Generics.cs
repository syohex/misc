global using IntGenerics = Generics<int>;

public class Generics<T>
{
    public T Field { get; set; }

    public Generics(T field) {
        Field = field;
    }
}
