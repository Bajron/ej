# ej

Eval JSON

# Usage

`ej <selector> [file]`

# Example

`ej Data.Status debug.json`

When no file is given, data is read from standard input.

`echo '[1,2,3]' | ej [1]`

# Why?

Sure, there exists jsawk or just python but if you simply want to stay in golang only ecosystem, this one might be helpful.

