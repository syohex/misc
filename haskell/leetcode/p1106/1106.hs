import Data.Bits

parseBoolExpr :: String -> Bool
parseBoolExpr s = fst $ parseBoolExpr' s

parseBoolExpr' :: String -> (Bool, String)
parseBoolExpr' [] = (True, [])
parseBoolExpr' s = case s of
    ('t' : xs) -> (True, xs)
    ('f' : xs) -> (False, xs)
    ('!' : '(' : xs) ->
        let (v, t) = parseBoolExpr' xs
         in (not v, tail t)
    ('&' : '(' : xs) ->
        let (exprs, t) = parseExprs xs []
         in (foldl (.&.) True exprs, t)
    ('|' : '(' : xs) ->
        let (exprs, t) = parseExprs xs []
         in (foldl (.|.) False exprs, t)
    _ -> error "never reach here"

parseExprs :: String -> [Bool] -> ([Bool], String)
parseExprs [] exprs = (exprs, [])
parseExprs s exprs =
    let (v, t) = parseBoolExpr' s
     in case t of
            (',' : t') -> parseExprs t' (v : exprs)
            (')' : t') -> (v : exprs, t')
            _ -> error "never reach here"
