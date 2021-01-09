# Bad Image (`.bi`)

`bi` is a library for the BI image format.
The BI format encodes images so that they are essentially CSV files using the named colors in the [CSS Color Module Level 4](https://www.w3.org/TR/css-color-4/#named-colors) or hex quadruplets to create a 2d matrix of color names.
There is also support for third party color models.
This means that they are plaintext editable.
Nice.

For example, a valid 2x2 pixel image is:
```
bi,cssColModLvl4
red,green
blue,white
```

## Examples

---

![Before BI encoding.](./misc/lenna1.png)

Before BI encoding.

---

![After BI encoding using the CSS Color Module Level 4 color model. And re-encoded as PNG obvs.](./misc/lenna2.png)

*After BI encoding using the CSS Color Module Level 4 color model. And re-encoded as PNG obvs.*

---

## License

This project is licensed under the MIT license.
