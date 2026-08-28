#!/usr/bin/env bun

import { say } from "../utils/utils";

async function main() {
    if (Bun.argv.length > 2) {
        say(Bun.argv.slice(2).join(", "))
        return;
    }
    const prompt = "Type something: ";
    process.stdout.write(prompt);
    for await (const line of console) {
        say(line);
        return
    }
}

main()
