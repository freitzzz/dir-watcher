# dir-watcher

My automation for making sure the `Downloads` folder doesn't get messy.

## What does it do?

It watches one or more directory for file system changes, moving old and new files to other directories where 
these files feel more appropriate in. For example:

> "family-photo.jpeg" > /home/user/Pictures
> "important-paper.pdf > /home/user/Documents/PDF

## Rules

The program uses a JSON file to know where to move files to. This JSON uses the example schema:

```json
{
    "watch": [
        "~/Downloads"
    ],
    "move": [
        {
            "path": "~/Pictures",
            "ext": [
                "png",
                "webp",
                "svg",
                "jpeg",
                "jpg",
                "tff"
            ]
        }
    ],
    "unknown": "~/Misc"
}
```

An extensive example can be seen in [rules.json](rules.json).