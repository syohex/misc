#[derive(Debug, PartialEq, Eq)]
pub enum Token {
    Integer(i64),
    Symbol(String),
    LParen,
    RParen,
}

impl std::fmt::Display for Token {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Token::Integer(n) => write!(f, "{n}"),
            Token::Symbol(s) => write!(f, "{s}"),
            Token::LParen => write!(f, "("),
            Token::RParen => write!(f, ")"),
        }
    }
}

pub fn tokenize(input: &str) -> Result<Vec<Token>, &str> {
    let mut tokens = Vec::new();
    let program = input.replace("(", " ( ").replace(")", " ) ");
    let words = program.split_whitespace();
    for word in words {
        match word {
            "(" => tokens.push(Token::LParen),
            ")" => tokens.push(Token::RParen),
            _ => {
                let v = word.parse::<i64>();
                if v.is_ok() {
                    tokens.push(Token::Integer(v.unwrap()));
                } else {
                    tokens.push(Token::Symbol(word.to_string()));
                }
            }
        }
    }

    Ok(tokens)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn add_expr() {
        let tokens = tokenize("(+ 1 2)").unwrap_or(vec![]);
        let expected = vec![
            Token::LParen,
            Token::Symbol("+".to_string()),
            Token::Integer(1),
            Token::Integer(2),
            Token::RParen,
        ];
        assert_eq!(tokens, expected);
    }
}
