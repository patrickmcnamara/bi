# Bad Image (`.bi`/`.biz`)

`bi` is a library for the `BI`/`BIZ` image format.
The `BI` format encodes images so that they are essentially CSV files using the named colors in the [CSS Color Module Level 4](https://www.w3.org/TR/css-color-4/#named-colors) to create a 2d matrix of color names.
This means that they are plaintext editable.
Nice.

For example, a valid 2x2 pixel image is:
```
bi
red, green
blue, white
```

The `BIZ` format is the same except it is gzipped and the magic number at the start is `biz\n` instead of `bi\n`.

## Examples

![Before BI encoding.](./misc/lenna1.png)

*Before BI encoding.*

![After BI encoding. And re-encoded as PNG obvs.](./misc/lenna2.png)

*After BI encoding. And re-encoded as PNG obvs.*

## License

This project is licensed under the GNU Lesser General Public License v3.0.
