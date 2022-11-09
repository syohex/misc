# C++ Software Design

## 読了日

## URL

https://www.oreilly.com/library/view/c-software-design/9781098113155/

## 感想

### Guideline 01 Understand the Importance of Software Design

- この本は「ソフトウェアデザイン」の本であって C++の機能の本ではないということを何度も強調

### Guideline 02 Design for Change

- 関心の分離, 単一責任原則, DRYについて
- 基本的な解答はない. すべて `It depends` だ
- 早すぎる関心の分離はやめておけ
- `Don't try to achieve SOLID, use SOLID to archive maintainability`

### Guideline 04 Design for Testability

- private methodのテストをどうするかという話
- protectedにしてテスト用のクラスを継承するとか, publicにするとか, friend使うとかマクロで privateを publicに置き換えるとか色々よろしくない方法もあげられる.
- publicメソッドを経由して private methodの機能を確認すればいいんじゃないという意見も出てくるが, 実装変わったりしたら当初のテストの目的を達成できないケースもある
- 最終的に提案されるものは, privateをやめて public関数に切り出すという方法. 引数に必要最小限なものを渡すようにすれば, classの中にあるより良いカプセル化もできるとのこと
- テストすべき機能はそれ自身テストできるようにあるべきだからそれでいいんじゃないとの意見. 同一クラスであることにこだわることは意味があることではないとの意見.

感想 単体テストであれば Rustみたいにファイル内でテストが書ければ別にいい話で言語の性質に引っ張られてテストが書けなかったり, 手間が必要だったりってのは今どきイケていないなと思った. テストするための設計の重要性はもちろんあるが, 余計なことを考えないように済むために言語規格がなんとかなってほしいと願うところではある

### Guideline 6 Adhere to the expected behavior of abstractions

- リスコフの置換原則の話
- 抽象化は requirementsと expectationsの表現である
- 数学的などでは成り立つ関係がソフトウェア的になりたつとは限らない. 例長方形と正方形. インターフェス次第ではリスコフの置換原則を満たさないという話を例を挙げて説明
