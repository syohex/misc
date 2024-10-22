let rec parse_bool_expr' (cs: char list) =
  match cs with
  | [] -> true, []
  | h :: t ->
     match h with
     | 't' -> true, t
     | 'f' -> false, t
     | '!' ->
        let v, t = parse_bool_expr' (List.tl t) in
        not v, List.tl t
     | '&' ->
        let v, t = parse_exprs (List.tl t) [] in
        List.fold_left (fun acc n -> acc && n) true v, t
     | '|' ->
        let v, t = parse_exprs (List.tl t) [] in
        List.fold_left (fun acc n -> acc || n) false v, t
     | _ -> failwith "never reach here"

and parse_exprs (cs: char list) acc =
  match cs with
  | [] -> failwith "never reach here"
  | ',' :: t -> parse_exprs t acc
  | ')' :: t -> acc, t
  | _ ->
     let v, t = parse_bool_expr' cs in
     parse_exprs t (v :: acc)

let parse_bool_expr (s: string) : bool =
  parse_bool_expr' (s |> String.to_seq |> List.of_seq) |> fst

let test () =
  begin
    Printf.printf "ret1=%b\n" (parse_bool_expr "&(|(f))");
    Printf.printf "ret2=%b\n" (parse_bool_expr "|(f,f,f,t)");
    Printf.printf "ret3=%b\n" (parse_bool_expr "!(&(f,t))");
    Printf.printf "ret4=%b\n" (parse_bool_expr "!(&(&(!(&(f)),&(t),|(f,f,t)),&(t),&(t,t,f)))");
  end
