import Data.List (sort, sortBy)
import Data.Ord

data Tree = Leaf | Node Int Tree Tree
    deriving (Eq)

kthLargestLevelSum' :: [Tree] -> Int -> [Int] -> Int
kthLargestLevelSum' [] k acc =
    if k - 1 >= length acc then -1 else sortBy (comparing Data.Ord.Down) acc !! (k - 1)
kthLargestLevelSum' q k acc =
    let (nodes, sum) =
            foldl
                ( \(acc, sum) node -> case node of
                    Leaf -> error "never reach here"
                    (Node v left right) ->
                        (right : left : acc, sum + v)
                )
                ([], 0)
                q
     in let q' = filter (/= Leaf) nodes
         in kthLargestLevelSum' q' k (sum : acc)

kthLargestLevelSum :: Tree -> Int -> Int
kthLargestLevelSum root k = kthLargestLevelSum' [root] k []

test :: IO ()
test = do
    let tree1 =
            Node
                5
                ( Node
                    8
                    ( Node
                        2
                        (Node 4 Leaf Leaf)
                        (Node 6 Leaf Leaf)
                    )
                    (Node 1 Leaf Leaf)
                )
                (Node 9 (Node 3 Leaf Leaf) (Node 7 Leaf Leaf))
    putStrLn $ "ret=" ++ show (kthLargestLevelSum tree1 2)
    putStrLn $ "ret2=" ++ show (kthLargestLevelSum tree1 50)
