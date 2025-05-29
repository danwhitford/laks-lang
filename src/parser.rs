use crate::lexer::Token;

#[derive(PartialEq, Debug)]
pub enum Stmt {
    // Some more stuff
    ExprStmt(Expr),
}

#[derive(PartialEq, Debug)]
pub enum Expr {
    Lit(Value),
    BinOp(Operator, Box<Expr>, Box<Expr>),
}

#[derive(PartialEq, Debug)]
pub enum Operator {
    ADD,
    SUB,
    MUL,
    DIV,
}

#[derive(PartialEq, Debug)]
pub enum Value {
    IntVal(i64),
}

pub fn parse(tokens: Vec<Token>) -> Vec<Stmt> {
    let mut exprs = Vec::new();
    let mut iter = tokens.iter().peekable();

    loop {
        match iter.peek() {
            Some(t) => match t {
                Token::Int(_) => exprs.push(Stmt::ExprStmt(parse_expr(&mut iter))),
                _ => panic!("dunno start of expression {:?}", t),
            },
            None => return exprs,
        }
        match iter.next() {
            Some(t) => match t {
                Token::Semi => (),
                _ => panic!("wanted semicolon got {:?}", t),
            },
            None => panic!("wanted semicolon got EOF"),
        }
    }
}

fn parse_expr<'a>(iter: &mut std::iter::Peekable<impl Iterator<Item = &'a Token>>) -> Expr {
    let mut expr = parse_expr2(iter);

    loop {
        match iter.peek() {
            Some(token) => match token {
                Token::Plus | Token::Sub => {
                    let op = token_to_op(iter.next().unwrap());
                    let right = parse_expr2(iter);
                    expr = Expr::BinOp(op, Box::from(expr), Box::from(right));
                }
                _ => return expr,
            },
            None => return expr,
        }
    }
}

fn parse_expr2<'a>(iter: &mut std::iter::Peekable<impl Iterator<Item = &'a Token>>) -> Expr {
    let mut expr = parse_literal(iter.next());

    loop {
        match iter.peek() {
            Some(token) => match token {
                Token::Mult | Token::Div => {
                    let op = token_to_op(iter.next().unwrap());
                    let right = parse_literal(iter.next());
                    expr = Expr::BinOp(op, Box::from(expr), Box::from(right));
                }
                _ => return expr,
            },
            None => return expr,
        }
    }
}

fn parse_literal(token: Option<&Token>) -> Expr {
    return match token {
        Some(token) => match token {
            Token::Int(s) => Expr::Lit(Value::IntVal(s.parse().unwrap())),
            _ => panic!("not a literal '{:?}'", token),
        },
        None => panic!("wanted literal got nothing!"),
    };
}

fn token_to_op(t: &Token) -> Operator {
    return match t {
        Token::Plus => Operator::ADD,
        Token::Sub => Operator::SUB,
        Token::Mult => Operator::MUL,
        Token::Div => Operator::DIV,
        _ => panic!("is not an operator {:?}", t),
    };
}

#[cfg(test)]
mod tests {
    // Note this useful idiom: importing names from outer (for mod tests) scope.
    use super::*;

    #[test]
    fn test_simple() {
        let got = parse(Vec::from([Token::Int(String::from("12")), Token::Semi]));
        let want = Vec::from([Stmt::ExprStmt(Expr::Lit(Value::IntVal(12)))]);

        assert_eq!(want, got);
    }

    #[test]
    fn test_addup() {
        let got = parse(Vec::from([
            Token::Int(String::from("12")),
            Token::Plus,
            Token::Int(String::from("107")),
            Token::Semi,
        ]));
        let want = Vec::from([Stmt::ExprStmt(Expr::BinOp(
            Operator::ADD,
            Box::from(Expr::Lit(Value::IntVal(12))),
            Box::from(Expr::Lit(Value::IntVal(107))),
        ))]);

        assert_eq!(want, got);
    }

    #[test]
    fn test_precedence() {
        let got = parse(Vec::from([
            Token::Int(String::from("1")),
            Token::Plus,
            Token::Int(String::from("2")),
            Token::Mult,
            Token::Int(String::from("3")),
            Token::Semi,
        ]));
        let want = Vec::from([Stmt::ExprStmt(Expr::BinOp(
            Operator::ADD,
            Box::from(Expr::Lit(Value::IntVal(1))),
            Box::from(Expr::BinOp(
                Operator::MUL,
                Box::from(Expr::Lit(Value::IntVal(2))),
                Box::from(Expr::Lit(Value::IntVal(3))),
            )),
        ))]);

        assert_eq!(want, got);
    }
}
