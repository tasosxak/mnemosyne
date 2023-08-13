# Mnemosyne - Event Data Preprocessing for DejaVu Runtime Monitor


Mnemosyne is a Go project designed to preprocess event data for the DejaVu runtime monitor.

## Features

- **Event Data Preprocessing:** Mnemosyne prepares event data for use with the DejaVu runtime monitor, making it easier to detect and analyze runtime events.


### Prerequisites

- Go (version 1.21.0)
- goyacc
- Nex

### Build
```sh
nex -o internals/lexer.go internals/lexer.nex \
&& goyacc -o internals/parser.go internals/parser.y \
&& go build -o mnemosyne internals/parser.go internals/lexer.go
```

### Example

```sh
event car:
input int speed, int speedLimit;
output int deltaSpeed, bool underLimit;
deltaSpeed <-  speed - @speed[speed];
underLimit <- speed <= speedLimit;
end
event detect:
input bool objectDetected;
output int speed;
speed <- ite(#objectDetected and speed < 5 , 0, 5);
end
```

### Run
```sh
mnemosyne<test.mn
```

The tools generates a .py file, you can execute it running the command below:

```sh
python code.py
```
