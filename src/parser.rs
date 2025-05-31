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
    let mut stmts = Vec::new();
    let mut iter = tokens.iter().peekable();

    while let Some(t) = iter.peek() {
        let stmt = match t {
            Token::Int(_) => Stmt::ExprStmt(parse_expr(&mut iter)),
            _ => panic!("dunno start of expression {:?}", t),
        };
        
        match iter.next() {
            Some(Token::Semi) => (),
            Some(t) => panic!("wanted semi got '{:?}'", t),
            None => panic!("wanted semi got EOF"),
        };

        stmts.push(stmt);
    }

    stmts
}

fn parse_expr<'a>(iter: &mut std::iter::Peekable<impl Iterator<Item = &'a Token>>) -> Expr {
    let mut expr = parse_expr2(iter);

    while let Some(token) = iter.next_if(|t| matches!(t, Token::Plus | Token::Sub)) {
        match token {
            Token::Plus | Token::Sub => {
                let op = token_to_op(token);
                let right = parse_expr2(iter);
                expr = Expr::BinOp(op, Box::from(expr), Box::from(right));
            }
            _ => panic!("shouldn't happen"),
        }
    }

    expr
}

fn parse_expr2<'a>(iter: &mut std::iter::Peekable<impl Iterator<Item = &'a Token>>) -> Expr {
    let mut expr = parse_literal(iter);

    while let Some(token) = iter.next_if(|t| matches!(t, Token::Mult | Token::Div)) {
        match token {
            Token::Mult | Token::Div => {
                let op = token_to_op(token);
                let right = parse_literal(iter);
                expr = Expr::BinOp(op, Box::from(expr), Box::from(right));
            }
            _ => panic!("shouldn't happen"),
        }
    }

    expr
}

fn parse_literal<'a>(iter: &mut std::iter::Peekable<impl Iterator<Item = &'a Token>>) -> Expr {
    match iter.next() {
        Some(Token::Int(s))=> Expr::Lit(Value::IntVal(s.parse().unwrap())),
        Some(token) =>  panic!("not a literal '{:?}'", token),
        None => panic!("wanted literal got EOF"),
    }
}

fn token_to_op(t: &Token) -> Operator {
    match t {
        Token::Plus => Operator::ADD,
        Token::Sub => Operator::SUB,
        Token::Mult => Operator::MUL,
        Token::Div => Operator::DIV,
        _ => panic!("is not an operator {:?}", t),
    }
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
