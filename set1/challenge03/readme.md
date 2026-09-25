## Overview

This challenge introduces **single-byte XOR encryption** and a basic form of **cryptanalysis**.

Unlike the previous challenge, where two known buffers were XORed together, here we are given a ciphertext encrypted with an **unknown single-byte key**.

The goal is to determine which key was used and recover the plaintext.

The challenge provides a ciphertext in hexadecimal form:

```text
1b37373331363f78151b7f2b783431333d78397828372d36
```

The important part of this challenge is not simply performing XOR. The interesting problem is:

> **How can we determine which of the 256 possible keys produces meaningful English text?**

---

# 1. XOR Refresher

XOR, or **exclusive OR**, operates on individual bits.

| A | B | A XOR B |
| - | - | ------- |
| 0 | 0 | 0       |
| 0 | 1 | 1       |
| 1 | 0 | 1       |
| 1 | 1 | 0       |

The important property for this challenge is that XOR is reversible:

```text
A XOR B XOR B = A
```

For example:

```text
1010 XOR 1100 = 0110
```

Then:

```text
0110 XOR 1100 = 1010
```

The original value is recovered.

This means that if encryption is:

```text
plaintext XOR key = ciphertext
```

then decryption is:

```text
ciphertext XOR key = plaintext
```

The same XOR operation is used in both directions.

---

# 2. What Is a Single-Byte XOR Cipher?

A **single-byte XOR cipher** uses one byte as the encryption key.

For example:

```text
Key = 0x42
```

A byte contains 8 bits, so the same key is XORed against every byte of the plaintext:

```text
plaintext[0] XOR 0x42
plaintext[1] XOR 0x42
plaintext[2] XOR 0x42
plaintext[3] XOR 0x42
...
```

The key does not change.

Conceptually:

```text
              same key
                  ↓
Plaintext → XOR → Ciphertext
```

For decryption:

```text
              same key
                  ↓
Ciphertext → XOR → Plaintext
```

---

# 3. Why Is It Called "Single-Byte"?

A byte contains 8 bits.

Therefore, there are:

```text
2^8 = 256
```

possible byte values.

The possible keys range from:

```text
0x00
```

through:

```text
0xFF
```

That gives us only **256 possible keys**.

This means we can try every possible key.

This technique is called **brute-force key search**.

---

# 4. The Basic Attack

We don't know the key.

So we can try:

```text
0x00
0x01
0x02
...
0xFF
```

For each key:

```text
ciphertext XOR key
```

produces a possible plaintext.

Conceptually:

```text
                 Ciphertext
                     │
       ┌─────────────┼─────────────┐
       ↓             ↓             ↓
     Key 0x00      Key 0x01      Key 0x02
       ↓             ↓             ↓
      XOR           XOR           XOR
       ↓             ↓             ↓
 Candidate 1    Candidate 2    Candidate 3
```

Continue until all 256 keys have been tested.

The problem is now:

> How do we identify which candidate is actual English?

---

# 5. Doing It by Hand: Use Common English Characters

English text contains characters with different frequencies.

Some very common characters are:

```text
space
e
t
a
o
i
n
s
h
r
```

The **space character** is particularly useful because spaces occur very frequently in normal English text.

ASCII space is:

```text
0x20
```

Suppose we think a ciphertext byte represents a space.

We know:

```text
ciphertext = plaintext XOR key
```

Therefore:

```text
key = ciphertext XOR plaintext
```

If the ciphertext byte is:

```text
0x78
```

and we guess that the plaintext byte is a space:

```text
0x20
```

then:

```text
0x78 XOR 0x20 = 0x58
```

So:

```text
candidate key = 0x58
```

We can now test this candidate key against the entire ciphertext.

This is an example of **known-plaintext reasoning**. We don't actually know the plaintext, but we are making an educated guess about one character based on the statistical properties of English.

---

# 6. Testing a Candidate Key

The ciphertext is:

```text
1b 37 37 33 31 36 3f 78 15 1b 7f 2b 78 34 31 33
3d 78 39 78 28 37 2d 36
```

Suppose we are testing the candidate key:

```text
0x58
```

XOR every byte with the same key.

The beginning becomes:

```text
1b XOR 58 = 43 = C
37 XOR 58 = 6f = o
37 XOR 58 = 6f = o
33 XOR 58 = 6b = k
31 XOR 58 = 69 = i
36 XOR 58 = 6e = n
3f XOR 58 = 67 = g
78 XOR 58 = 20 = space
```

The resulting beginning is therefore:

```text
Cooking 
```

This is a strong indication that this candidate key is worth investigating further because the output has immediately started producing recognizable English.

Continue applying the same key to the remaining ciphertext bytes and evaluate the complete result.

> **Do not assume a candidate is correct just because the first few characters look meaningful.** The complete plaintext should be coherent English.

---

# 7. Why Repeated Ciphertext Bytes Are Useful

The same key is used for every byte.

Suppose two ciphertext positions contain the same value:

```text
37
37
```

Because:

```text
ciphertext = plaintext XOR key
```

the corresponding plaintext characters must also be the same.

For example:

```text
P1 XOR K = 37
P2 XOR K = 37
```

Therefore:

```text
P1 = P2
```

Repeated ciphertext bytes can therefore provide clues about repeated plaintext characters.

This becomes especially useful when combined with guesses about spaces and common letters.

---

# 8. Why Scoring Is Necessary

Manually looking at candidate plaintexts works for a small challenge, but it does not scale well.

Instead, we can give every candidate plaintext a **score**.

The score represents:

> How closely does this candidate resemble English text?

For example, we could assign positive points to common characters:

```text
space → high score
e     → high score
t     → high score
a     → high score
o     → high score
i     → high score
n     → high score
```

Less common characters could receive smaller scores.

Non-printable characters could receive a strong penalty.

For example:

```text
English-looking text:
"the quick brown..."
→ high score

Random binary:
"\x01\x9a\x03..."
→ very low score
```

The exact numerical values are not important.

What matters is that **English-like output should generally score higher than random output**.

---

# 9. Simple Character-Frequency Scoring

A basic scoring model could look conceptually like:

```text
Character              Score

space                    +5
e                        +4
t                        +4
a                        +3
o                        +3
i                        +3
n                        +3
s                        +2
h                        +2
r                        +2

other printable          +1

non-printable            -10
```

These numbers are only an example.

The scoring system does not need to perfectly understand English.

It only needs to make likely English plaintext score higher than random-looking output.

---

# 10. Example of Scoring

Suppose two candidate plaintexts are:

```text
Candidate A:
the quick brown
```

and:

```text
Candidate B:
\x01\x9F@#\x03\x11
```

Candidate A contains:

* spaces
* common letters
* printable characters
* recognizable English patterns

Candidate B contains:

* non-printable bytes
* unusual characters
* no spaces
* no recognizable English structure

Therefore:

```text
Candidate A → high score
Candidate B → low score
```

The candidate with the highest score becomes our best candidate.

---

# 11. Character Frequency Is Not Perfect

A simple character-frequency system has limitations.

For example:

```text
eeeeeeeeeeeeeeeeeeee
```

contains a very common English character but is not realistic English.

More sophisticated scoring can consider:

### Character frequency

How often individual characters appear.

### Bigrams

Common two-character combinations:

```text
th
he
in
er
an
re
```

### Trigrams

Common three-character combinations:

```text
the
and
ing
```

### Common words

For example:

```text
the
and
that
this
with
from
```

### Printable characters

Normal English text should mostly contain printable ASCII characters.

For this challenge, a simple scoring system is sufficient.

---

# 12. Automated Brute Force

Once we have a scoring system, the attack becomes straightforward.

```text
Ciphertext
    │
    ↓
Try key 0x00
    │
    ↓
Decrypt
    │
    ↓
Score plaintext
    │
    ↓
Try key 0x01
    │
    ↓
Decrypt
    │
    ↓
Score plaintext
    │
    ↓
...
    │
    ↓
Try key 0xFF
    │
    ↓
Decrypt
    │
    ↓
Score plaintext
    │
    ↓
Select highest-scoring candidate
```

The computer is essentially automating the manual process.

Instead of asking us:

> "Does this look like English?"

for every key, we give the computer a scoring function that approximates that judgment.

---

# 13. Important Formulas

The encryption relationship is:

```text
ciphertext = plaintext XOR key
```

Because XOR is reversible:

```text
plaintext = ciphertext XOR key
```

And if both plaintext and ciphertext are known:

```text
key = plaintext XOR ciphertext
```

These relationships are fundamental to understanding XOR-based cryptanalysis.

---

# 14. Challenge Workflow

The complete workflow is:

```text
Hexadecimal ciphertext
        │
        ↓
Decode hexadecimal
        │
        ↓
Byte buffer
        │
        ↓
Try all 256 possible keys
        │
        ↓
XOR ciphertext with candidate key
        │
        ↓
Candidate plaintext
        │
        ↓
Score candidate
        │
        ↓
Keep highest-scoring result
```

---

