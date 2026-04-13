import cowsay from "cowsay";

export function say(text: string)  {
    console.log(cowsay.say({ text }))
}
