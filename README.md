# Spanish numbers

Using https://langpractice.com to provide a command-line utility.

Motivation for this instead of the webapp from langpractice.com is to optionally wait on user before asking for the next random number sample.

## Usage
```shell
Usage:
  -max int
        max (default 1000)
  -min int
        min (default 100)
  -noan
        no auto next number
```

## LangPractice API
### Request
URL: https://langpractice.com/api/newnumber/spanish-mexico/1/100/1000 - ".../1/min/max"

### Response
```json
{
  "audio_data": "...base64-encoded MPEG ADTS, layer III, v2, 160 kbps, 24 kHz, Monaural",
  "n": 380,
  "target": {
    "phonetic": "trescientos ochenta",
    "written": "trescientos ochenta"
  }
}
```