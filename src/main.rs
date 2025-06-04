use std::io::Write;

use laks_lang::{binterp, compiler::compile, lexer::lex, parser::parse};


fn run(source: &str, out: &mut impl Write) {
    let tokens = lex(source);
    let stmts = parse(tokens);
    let chunk = compile(stmts);
    binterp::run(chunk, out);
}

fn main() {
    println!("Hello, world!");
}

#[cfg(test)]
mod tests {
    // Note this useful idiom: importing names from outer (for mod tests) scope.
    use super::*;

    #[test]
    fn test_intval() {
        let input = "print 107;";

        let want = "107\n";

        let mut buf = vec![];
        run(input, &mut buf);
        let got = String::from_utf8(buf).expect("should be valid utf8 string");

        assert_eq!(want, got);
    }
}
