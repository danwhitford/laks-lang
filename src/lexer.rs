
struct Lexer {
    source: String,
}

#[derive(PartialEq, Debug)]
enum Token {
    TokenInt(String),
    TokenPlus,
    TokenSub,
    TokenMult,
    TokenDiv,
}

impl Lexer {
    fn lex(&self) -> Vec<Token> {
        let mut tots = Vec::new();
        let mut foo = self.source.chars().peekable();

        loop {
            match foo.next() {
                Some(c) => match c {
                    '0'..'9' => tots.push(self.number(c,&mut foo)),
                    '+' => tots.push(Token::TokenPlus),
                    '-' => tots.push(Token::TokenSub),
                    '*' => tots.push(Token::TokenMult),
                    '/' => tots.push(Token::TokenDiv),
                    c if c.is_whitespace() => (),
                    _ => panic!("don't know this character '{c}'"),
                },
                None => return tots,
            }
        }
    }

    fn number(&self, c: char, iter: &mut impl Iterator<Item = char>) -> Token {
      let mut lexeme = String::from(c);
      loop {
        match iter.next() {
          Some(c) => match c {
            '0'..'9' => lexeme.push(c),
            c if c.is_whitespace() => return Token::TokenInt(lexeme),
            _ => panic!("not part of a number '{c}'"),
          }
          None => return Token::TokenInt(lexeme)
        }
      }
    }
}

#[cfg(test)]
mod tests {
    // Note this useful idiom: importing names from outer (for mod tests) scope.
    use super::*;

    #[test]
    fn test_simple() {
        let l = Lexer {
            source: String::from("4"),
        };
        let mut want = Vec::new();
        want.push(Token::TokenInt(String::from("4")));

        assert_eq!(want, l.lex());
    }

    #[test]
    fn test_nums() {
        let l = Lexer {
            source: String::from("1 2 3"),
        };
        let want = Vec::from([
            Token::TokenInt(String::from("1")),
            Token::TokenInt(String::from("2")),
            Token::TokenInt(String::from("3")),
        ]);

        assert_eq!(want, l.lex());
    }

    #[test]
    fn test_real_big_num() {
        let l = Lexer {
            source: String::from("478"),
        };
        let want = Vec::from([
            Token::TokenInt(String::from("478")),
        ]);

        assert_eq!(want, l.lex());
    }

    #[test]
    fn basic_operators() {
      let l = Lexer {
        source: String::from("+ - / *"),
    };
    let want = Vec::from([
        Token::TokenPlus,
        Token::TokenSub,
        Token::TokenDiv,
        Token::TokenMult,
    ]);

    assert_eq!(want, l.lex());
    }
}
