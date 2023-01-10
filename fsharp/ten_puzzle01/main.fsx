type Operation =
    | Add
    | Sub
    | Mul
    | Div

type Node =
    | Operation of Operation
    | Value of int

let calcRPN (rpn: Node list) : Result<int, string> =
    let rec calcRPN' rpn stack =
        match rpn with
        | [] -> Ok(List.head stack)
        | Value (v) :: t -> calcRPN' t (v :: stack)
        | Operation (op) :: t ->
            match stack with
            | v2 :: v1 :: rest ->
                match op with
                | Add -> calcRPN' t ((v1 + v2) :: rest)
                | Sub -> calcRPN' t ((v1 - v2) :: rest)
                | Mul -> calcRPN' t ((v1 * v2) :: rest)
                | Div ->
                    if v2 = 0 then
                        Error("zero divide")
                    else
                        calcRPN' t ((v1 / v2) :: rest)
            | _ -> failwith "invalid stack"

    calcRPN' rpn []

let tenPuzzle (nums: int list) : Node list list =
    let operations = [ Add; Sub; Mul; Div ]

    let rec tenPuzzle' (nums: int list) (values: int) (ops: int) limit stack acc =
        if values + ops >= limit then
            let stack' = stack |> List.rev

            match calcRPN stack' with
            | Ok (v) -> if v = 10 then stack' :: acc else acc
            | Error (_) -> acc
        else
            let acc' =
                if values - ops >= 2 then
                    operations
                    |> List.fold (fun acc op -> tenPuzzle' nums values (ops + 1) limit (Operation(op) :: stack) acc) acc
                else
                    acc

            match nums with
            | [] -> acc'
            | h :: t -> tenPuzzle' t (values + 1) ops limit (Value(h) :: stack) acc'

    let len = nums.Length
    let limit = len + (len - 1)
    tenPuzzle' nums 0 0 limit [] []

tenPuzzle [ 1; 3; 5; 8 ]
tenPuzzle [ 2; 3; 4; 5 ]
