type tree = Leaf | Node of int * tree * tree

let kth_largest_level_sum (root : tree) (k : int) =
  let rec f node level acc =
    match node with
    | Leaf -> acc
    | Node (v, left, right) ->
        let sum = Hashtbl.find_opt acc level in
        let acc' =
          match sum with
          | Some s ->
              Hashtbl.add acc level (s + v);
              acc
          | None ->
              Hashtbl.add acc level v;
              acc
        in
        let acc' = f left (level + 1) acc' in
        f right (level + 1) acc'
  in

  let sums = f root 0 (Hashtbl.create 1024) in
  let v = Hashtbl.to_seq_values sums |> List.of_seq in
  let v = List.sort compare v |> List.rev in
  let v = List.nth_opt v (k - 1) in
  match v with None -> -1 | Some v -> v

let test () =
  let tree1 =
    Node
      ( 5,
        Node
          ( 8,
            Node (2, Node (4, Leaf, Leaf), Node (6, Leaf, Leaf)),
            Node (1, Leaf, Leaf) ),
        Node (9, Node (3, Leaf, Leaf), Node (7, Leaf, Leaf)) )
  in
  Printf.printf "ret1=%d\n" (kth_largest_level_sum tree1 2);
  Printf.printf "ret2=%d\n" (kth_largest_level_sum tree1 50);
  let tree2 = Node (1, Node (2, Node (3, Leaf, Leaf), Leaf), Leaf) in
  Printf.printf "ret3=%d\n" (kth_largest_level_sum tree2 1)
