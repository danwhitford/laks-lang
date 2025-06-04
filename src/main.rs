use std::io::{stdout, Write};

use crate::{compiler::compile, lexer::lex, parser::parse};

mod binterp;
mod compiler;
mod lexer;
mod parser;

fn run(source: &str, out: &mut impl Write) {
    let tokens = lex(source);
    let stmts = parse(tokens);
    let chunk = compile(stmts);
    binterp::run(chunk, out);
}

fn main() {
    run("print 42;", &mut stdout());
}

#[cfg(test)]
mod tests {
    // Note this useful idiom: importing names from outer (for mod tests) scope.
    use super::*;

    fn test_case(source: &str, want: &str) {
        let mut buf = vec![];
        run(source, &mut buf);
        let got = String::from_utf8(buf).expect("should be valid utf8 string");

        assert_eq!(want, got, "test failed for '{}'", source);
    }

    #[test]
    fn test_intval() {
        let input = "print 107;";

        let want = "107\n";

        let mut buf = vec![];
        run(input, &mut buf);
        let got = String::from_utf8(buf).expect("should be valid utf8 string");

        assert_eq!(want, got);
    }

    #[test]
    fn table_tests() {
        test_case("print 42;", "42\n");
        test_case("print 2 + 2;", "4\n");
    }
}
