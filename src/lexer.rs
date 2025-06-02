#[derive(PartialEq, Debug)]
pub enum Token {
    Int(String),
    Plus,
    Sub,
    Mult,
    Div,
    Semi,
}

pub fn lex(source: &str) -> Vec<Token> {
    let mut tots = Vec::new();
    let mut chars = source.chars().peekable();

    while let Some(c) = chars.next() {
        let tok = match c {
            '0'..'9' => Some(number(&mut chars, c)),
            '+' => Some(Token::Plus),
            '-' => Some(Token::Sub),
            '*' => Some(Token::Mult),
            '/' => Some(Token::Div),
            ';' => Some(Token::Semi),
            c if c.is_whitespace() => None,
            _ => panic!("don't know this character '{c}'"),
        };

        match tok {
            Some(t) => tots.push(t),
            None => (),
        }
    }
    tots
}

fn number(iter: &mut std::iter::Peekable<impl Iterator<Item = char>>, c: char) -> Token {
    let mut lexeme = String::from(c);

    while let Some(c) = iter.next_if(|c| matches!(c, '0'..'9')) {
        match c {
            '0'..'9' => lexeme.push(c),
            _ => panic!("not part of a number '{c}'"),
        }
    }

    Token::Int(lexeme)
}

#[cfg(test)]
mod tests {
    // Note this useful idiom: importing names from outer (for mod tests) scope.
    use super::*;

    #[test]
    fn test_simple() {
        let l = lex("4");
        let mut want = Vec::new();
        want.push(Token::Int(String::from("4")));

        assert_eq!(want, l);
    }

    #[test]
    fn test_nums() {
        let l = lex("1 2 3");
        let want = Vec::from([
            Token::Int(String::from("1")),
            Token::Int(String::from("2")),
            Token::Int(String::from("3")),
        ]);

        assert_eq!(want, l);
    }

    #[test]
    fn test_real_big_num() {
        let l = lex("478");
        let want = Vec::from([Token::Int(String::from("478"))]);

        assert_eq!(want, l);
    }

    #[test]
    fn basic_operators() {
        let l = lex("+ - / *");
        let want = Vec::from([Token::Plus, Token::Sub, Token::Div, Token::Mult]);

        assert_eq!(want, l);
    }

    #[test]
    fn full_expr() {
        let l = lex("8 + 4;");
        let want = Vec::from([
            Token::Int(String::from("8")),
            Token::Plus,
            Token::Int(String::from("4")),
            Token::Semi,
        ]);

        assert_eq!(want, l);
    }
}
