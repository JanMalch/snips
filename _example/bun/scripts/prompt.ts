#!/usr/bin/env bun

import { say } from "../utils/utils";

async function main() {
    const prompt = "Type something: ";
    process.stdout.write(prompt);
    for await (const line of console) {
        say(line);
        return
    }
}

main()
