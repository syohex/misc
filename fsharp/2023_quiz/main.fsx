type Opcode =
    | Add
    | Sub
    | Mul
    | Div

type Operand =
    | Opcode of Opcode
    | Number of int

let calculateRpnExpression (exp: Operand list) : Option<int> =
    let rec calculateRpnExpression' exp stack =
        match exp with
        | [] -> Some(List.head stack)
        | Number (n) :: t -> calculateRpnExpression' t (n :: stack)
        | Opcode (op) :: t ->
            match stack with
            | n2 :: n1 :: rest ->
                match op with
                | Add -> calculateRpnExpression' t ((n1 + n2) :: rest)
                | Sub -> calculateRpnExpression' t ((n1 - n2) :: rest)
                | Mul -> calculateRpnExpression' t ((n1 * n2) :: rest)
                | Div ->
                    if n2 = 0 then
                        None
                    else
                        calculateRpnExpression' t ((n1 / n2) :: rest)
            | _ ->
                printfn $"stack=${stack} exp=${exp}"
                failwith "never reach here"

    calculateRpnExpression' exp []

let normalizeRpnExpression (exp: Operand list) : string =
    let rec normalizeRpnExpression' exp stack =
        match exp with
        | [] -> List.head stack
        | Number (n) :: t -> normalizeRpnExpression' t ((string n) :: stack)
        | Opcode (op) :: t ->
            match stack with
            | v2 :: v1 :: rest ->
                match op with
                | Add -> normalizeRpnExpression' t ((sprintf "(%s+%s)" v1 v2) :: rest)
                | Sub -> normalizeRpnExpression' t ((sprintf "(%s-%s)" v1 v2) :: rest)
                | Mul -> normalizeRpnExpression' t ((sprintf "(%s*%s)" v1 v2) :: rest)
                | Div -> normalizeRpnExpression' t ((sprintf "(%s/%s)" v1 v2) :: rest)
            | _ -> failwith "never reach here"

    normalizeRpnExpression' exp []

let createRpnExpression (nums: int list) (target: int) =
    let len = nums |> List.length
    let limit = len * 2 - 1

    let rec createRpnExpression' nums prevNum numCount opCount exp acc =
        if numCount + opCount >= limit then
            let exp' = List.rev exp

            match calculateRpnExpression exp' with
            | None -> acc
            | Some (v) ->
                if v = target then
                    printfn "%A" (normalizeRpnExpression exp')
                    exp' :: acc
                else
                    acc
        else
            let acc' =
                if numCount - opCount >= 2 then
                    let ops =
                        if prevNum = 8 then
                            [ Div ]
                        else
                            [ Add; Sub; Mul; Div ]

                    ops
                    |> List.fold
                        (fun acc' op ->
                            createRpnExpression' nums -1 numCount (opCount + 1) (Opcode(op) :: exp) acc')
                        acc
                else
                    acc

            if prevNum <> 8 && numCount < len then
                let n = List.head nums
                createRpnExpression' (List.tail nums) n (numCount + 1) opCount (Number(n) :: exp) acc'
            else
                acc'

    createRpnExpression' nums -1 0 0 [] []

let quiz2023 (nums: int list) (target: int) : string list =
    let rpnExpressions = createRpnExpression nums target
    rpnExpressions |> List.map normalizeRpnExpression

// quiz2023 [ 1; 2; 3; 4 ] 10

let input = seq { 1..10 } |> Seq.rev |> Seq.toList
quiz2023 input 2023 |> List.take 3
