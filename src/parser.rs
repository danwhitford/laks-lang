use crate::lexer::Token;

#[derive(PartialEq, Debug)]
enum Stmt {
    // Some more stuff
    ExprStmt(Expr),
}

#[derive(PartialEq, Debug)]
enum Expr {
    Lit(Value),
    BinOp(Operator, Box<Expr>, Box<Expr>),
}

#[derive(PartialEq, Debug)]
enum Operator {
    ADD,
    SUB,
    MUL,
    DIV,
}

#[derive(PartialEq, Debug)]
enum Value {
    IntVal(i64),
}

fn parse(tokens: Vec<Token>) -> Vec<Stmt> {
    let mut exprs = Vec::new();
    let mut iter = tokens.iter().peekable();

    loop {
        match iter.peek() {
            Some(t) => match t {
                Token::Int(_) => {
                    exprs.push(Stmt::ExprStmt(parse_expr(&mut iter)))
                }
                _ => panic!("dunno start of expression {:?}", t),
            },
            None => return exprs,
        }
    }
}

fn parse_expr<'a>(iter: &mut std::iter::Peekable<impl Iterator<Item = &'a Token>>) -> Expr {
    let mut expr = parse_literal(iter.next());

    loop {
        match iter.peek() {
            Some(token) => match token {
                Token::Plus | Token::Sub => {
                    let op = token_to_op(iter.next().unwrap());
                    let right = parse_literal(iter.next());
                    expr = Expr::BinOp(op, Box::from(expr), Box::from(right));
                }
                _ => (),
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

fn int_expr(s: &str) -> crate::parser::Value {
    return Value::IntVal(s.parse().unwrap());
}

#[cfg(test)]
mod tests {
    // Note this useful idiom: importing names from outer (for mod tests) scope.
    use super::*;

    #[test]
    fn test_simple() {
        let got = parse(Vec::from([Token::Int(String::from("12"))]));
        let want = Vec::from([Stmt::ExprStmt(Expr::Lit(Value::IntVal(12)))]);

        assert_eq!(want, got);
    }

    #[test]
    fn test_addup() {
        let got = parse(Vec::from([
            Token::Int(String::from("12")),
            Token::Plus,
            Token::Int(String::from("107")),
        ]));
        let want = Vec::from([Stmt::ExprStmt(Expr::BinOp(
            Operator::ADD,
            Box::from(Expr::Lit(Value::IntVal(12))),
            Box::from(Expr::Lit(Value::IntVal(107))),
        ))]);

        assert_eq!(want, got);
    }
}
