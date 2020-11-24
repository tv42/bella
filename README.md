# Bella -- label printer software

![`bel   la`](bella.png)

Bella renders text to graphics and prints it on a label maker using IPP/CUPS.

It was created for use with a [Dymo LabelManager PnP](https://smile.amazon.com/gp/product/B00464E5P2/), but other continuous-tape printers might work too.
If you have a different brand/model printer, let us know how to get it working!

**Current status: just getting started on the project.**

```
bella-png >bella.png 'bel   la'
lp -d label -o landscape bella.png
```

## Roadmap

- command-line printing
- web interface, especially for mobile client use
- submit print jobs via IPP?
- limit label maximum width (mm or inch), to e.g. fit on drawer front
- maybe embed images?
- maybe embed QR codes?

## Supported outputs and tapes

We currently assume the output needs to be 64 pixels tall.
This matches the printable area of a 1/2-inch Dymo D1 tape.

### Dymo D1 tape widths

| Tape size (inch) | Tape size (mm, approx) | Printable area (mm) | Printable area (pixels) |
|------------------|------------------------|---------------------|-------------------------|
| 1/2              | 12                     | 9                   | 64                      |
| 3/8              | 9                      | ?                   | ?                       |
| 1/4              | 6                      | ?                   | ?                       |


## Setting up CUPS for Dymo (on Debian/Ubuntu)

See https://www.baitando.com/it/2017/12/12/install-dymo-labelwriter-on-headless-linux

```
sudo apt-get install printer-driver-dymo

# for me, the above didn't include the PPD file for some reason; get it from the SDK
wget http://download.dymo.com/dymo/Software/Download%20Drivers/Linux/Download/dymo-cups-drivers-1.4.0.tar.gz
tar xzf dymo-cups-drivers-1.4.0.tar.gz
sudo cp dymo-cups-drivers-1.4.0.5/ppd/lmpnp.ppd /usr/share/cups/model/

# find your device serial number in lsusb, it's a 14-digit number
lsusb -d 0922:1002 -v
lpadmin -p label -v 'usb://DYMO/LabelManager%20PnP?serial=REPLACE_SERIAL_HERE' -P /usr/share/cups/model/lmpnp.ppd
cupsenable label

# test it
convert -size 300x64 canvas:white -font 'Times-Roman' -pointsize 64 -fill black -draw 'text 0,64 "label maker"' label.png
lp -d label -o landscape label.png
```


## Open questions

### How many levels of gray can a Dymo print?

We're currently using 8-bit grayscale. Is there any point in using 16-bit?

### Automatic cutting

Do label printers with automatic cutting need something special?
Right now we just print it out, cutting is end user problem (and the Dymo LabelManager PnP has a manual cutter); this means if you print multiple labels in a row, you lose the use of the built-in cutter and need scissors.
You can configure the Dymo LabelManager PnP to pause after each label, to leave you time to cut, but I haven't tested if it actually waits for user action or not.


## Anti-goals

- Fixed size pre-cut labels. We'd rather get the continuous tape style labels working *great*, and that means letting the length of the printout be whatever it is.
- Generating PDF. The printer driver will just rasterize it anyway, so we can simply submit a bitmap. The size isn't an issue here.


## Alternatives

- You could use the Windows and macOS software that came with the label maker.
- You could make bitmaps with ImageMagick: https://unix.stackexchange.com/questions/138804/how-to-transform-a-text-file-into-a-picture


## Resources

- You'll need the Dymo SDK to get the `PPD` file: https://www.dymo.com/en-US/dymo-label-sdk-cups-linux-p (that page is horribly broken, but the "Download" link works).
  Announcement at <https://developers.dymo.com/2012/02/21/announcing-dymo-labelwriterlabelmanager-sdk-1-4-0-for-linux/>.
  Fluffy landing page at <https://www.dymo.com/en-US/online-support/online-support-sdk>.
  Unofficial mirrors at https://github.com/matthiasbock/dymo-cups-drivers and https://github.com/Kyle-Falconer/DYMO-SDK-for-Linux and https://github.com/xcross/dymo-cups-drivers

### Random related things, not endorsing just cataloging

- https://www.linux-magazine.com/Issues/2016/183/Perl-Producing-Labels
- The HID device and UBS mode switch stuff seems to be handled fine by `printer-driver-dymo` these days, you no longer need the reverse-engineered `dymoprint`:
  https://sbronner.com/dymoprint.html
  https://github.com/computerlyrik/dymoprint
  https://randomfoo.net/2018/07/09/printing-with-a-dymo-labelmanager-pnp-on-linux
  or this alternate implementation https://github.com/Firedrake/dymo-labelmanager
- https://glabels.org/ or https://github.com/jimevins/glabels-qt seems more for pre-cut stickers ([mostly](https://github.com/jimevins/glabels-qt/commit/467ca9fc624e07442d45b0214f2cccb4919004a4))
- https://kwagjj.wordpress.com/2017/05/10/dymo-linux-command-line/ talks more about gLabels and its batch mode, which I do not wish to use.
- Something in C# https://github.com/ChaliceAriel/LabelManager
- Something in Python https://github.com/richardbarlow/hacky-dymo-label-gen
- dymoprint wrapper https://github.com/matrach/iot-labeler that seems to contain no actual code
