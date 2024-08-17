# 🎉 truvami Decoder 🚀

**truvami Decoder** is a command-line interface (CLI) tool written in Go for decoding truvami payloads. This reference implementation supports various payload types, including Nomad XS and different Tag formats. 🛠️

## 🎬 Demo

Check out this demo to see truvami Decoder in action! 👇

![Demo](demo.gif)

## 🌟 Features

- 🔍 **Payload Decoding**: Decode Nomad XS, Tag S/L, and Tag XL payloads with ease.
- 📄 **Flexible Output**: Choose between standard console output and JSON format.
- 🛠️ **Debugging & Verbosity**: Enable debugging and verbose output for detailed insights.

## 📦 Installation

### 🔧 Binaries

You can install the truvami Decoder binary easily using the following command:

```zsh
# This will install the binary at $(go env GOPATH)/bin/decoder
curl -sSfL https://raw.githubusercontent.com/truvami/decoder/main/install.sh | sh -s -- -b $(go env GOPATH)/bin

# ✅ Verify the installation by checking the help
decoder --help
```

### 🖥️ Windows

For Windows users, you can install the truvami Decoder binary using Chocolatey:

```powershell
# Install using Chocolatey
choco install truvami-decoder

# ✅ Verify the installation by checking the help
decoder --help
```

## 🛠️ Usage

truvami Decoder provides a variety of commands and options to help you decode payloads efficiently. Below is an overview of the available commands and flags.

### 🗂️ Basic Command Structure

```sh
decoder [command] [flags]
```

### 📝 Available Commands

- `completion` - 🖋️ Generate the autocompletion script for the specified shell.
- `help` - ℹ️ Display help information about any command.
- `nomadxs` - 🧩 Decode Nomad XS payloads.
- `tagsl` - 🏷️ Decode Tag S / L payloads.
- `tagxl` - 🏷️ Decode Tag XL payloads.

### 🚩 Global Flags

- `-d, --debug` - 🐛 Display debugging output in the console. (default: false)
- `-h, --help` - ℹ️ Display help information.
- `-j, --json` - 📄 Output the result in JSON format. (default: false)
- `-v, --verbose` - 📢 Display more verbose output in the console. (default: false)

### 💡 Example Usage

```sh
# 🔍 Decode a Nomad XS payload with verbose output
decoder nomadxs 1 0002c420ff005ed85a12b4180719142607240001ffbaffc2fc6f09a71d2e

# 📝 Decode a Tag S / L payload and output the result in JSON format
decoder tagsl 1 8002cdcd1300744f5e166018040b14341a -j

# 🖋️ Generate autocompletion script for bash
decoder completion bash
```

For more detailed information on each command, use:

```sh
decoder [command] --help
```
