IFoo f1 = new Foo1();
IFoo f2 = new Foo2();

Console.WriteLine($"f1 type={f1.GetType().FullName}");
Console.WriteLine($"f2 type={f2.GetType().FullName}");

interface IFoo
{
}

class Foo1 : IFoo
{
}

class Foo2 : IFoo
{
}

