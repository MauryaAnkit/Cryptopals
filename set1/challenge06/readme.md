# Challenge 6 — Break Repeating-Key XOR

## Objective

Break a ciphertext encrypted using a **repeating-key XOR cipher** when the key is unknown.

Unlike Challenge 3, where a single byte is reused as the key, repeating-key XOR uses a multi-byte key that repeats across the plaintext.

The challenge is to:

1. Determine the likely key size.
2. Recover each byte of the key.
3. Reconstruct the complete key.
4. Decrypt the ciphertext.

---

## How Repeating-Key XOR Works

Suppose the key is:

```text
ICE
```

The key bytes repeat continuously:

```text
Plaintext:  P0 P1 P2 P3 P4 P5 P6 P7 P8 ...
Key:        I  C  E  I  C  E  I  C  E  ...
```

Each plaintext byte is XORed with the corresponding key byte:

```text
Ciphertext[i] = Plaintext[i] XOR Key[i % KeySize]
```

The `%` operator makes the key repeat.

For a key of size `3`:

```text
i       i % 3
0       0
1       1
2       2
3       0
4       1
5       2
...
```

So:

```text
Key index:
0 1 2 0 1 2 0 1 2 ...
I C E I C E I C E ...
```

---

# The Problem

We have the ciphertext, but we don't know:

```text
Key size
Key
Plaintext
```

We need to recover them.

The basic attack is:

```text
Ciphertext
    │
    ▼
Try key sizes 2–40
    │
    ▼
Calculate normalized Hamming distance
    │
    ▼
Select promising key sizes
    │
    ▼
Transpose ciphertext
    │
    ▼
Each group becomes single-byte XOR
    │
    ▼
Use Challenge 3 technique
    │
    ▼
Recover each key byte
    │
    ▼
Reconstruct key
    │
    ▼
Decrypt ciphertext
```

---

# Step 1 — Guess the Key Size

The first problem is determining how long the key is.

For example, if:

```text
Key = ICE
```

then:

```text
Key size = 3
```

We can test possible key sizes, for example:

```text
2, 3, 4, 5, ... 40
```

For each possible key size, we calculate the **Hamming distance** between ciphertext blocks.

---

# Hamming Distance

Hamming distance measures how many **bits differ** between two equal-length byte sequences.

For example:

```text
A = 0x48
B = 0x6A
```

Convert them to binary:

```text
0x48 = 01001000
0x6A = 01101010
```

XOR them:

```text
01001000
01101010
--------
00100010
```

There are two `1` bits.

Therefore:

```text
Hamming Distance = 2
```

In Go, this can be calculated using:

```go
bits.OnesCount8(a ^ b)
```

---

# Why Use Hamming Distance?

If the key size is correct, ciphertext blocks tend to have a statistical relationship because the same key pattern is being reused.

For example, with:

```text
Key = ICE
```

and a key size of `3`:

```text
Block 1 → encrypted using ICE
Block 2 → encrypted using ICE
Block 3 → encrypted using ICE
```

The blocks are aligned with the repeating key.

For an incorrect key size, the key alignment shifts between blocks.

Therefore, comparing blocks using Hamming distance can help identify likely key sizes.

However, this is only a **heuristic**. The lowest distance is not guaranteed to be the correct key size.

---

# Normalized Hamming Distance

Different key sizes produce blocks containing different numbers of bits.

For example:

```text
Key size = 2
Maximum distance = 16 bits

Key size = 4
Maximum distance = 32 bits
```

Therefore, we cannot directly compare the raw distances.

Instead, normalize the distance:

```text
Normalized Distance =
    Hamming Distance / Key Size
```

For example:

```text
Key Size = 3
Hamming Distance = 10

Normalized Distance = 10 / 3
                    = 3.33
```

A **smaller normalized distance** makes a key size more interesting as a candidate.

---

# Comparing Multiple Blocks

We don't have to compare only two blocks.

For example, we can divide the ciphertext into blocks:

```text
Block 0
Block 1
Block 2
Block 3
Block 4
```

Then calculate:

```text
Distance(Block 0, Block 1)
Distance(Block 1, Block 2)
Distance(Block 2, Block 3)
Distance(Block 3, Block 4)
```

Then calculate the average.

The number of blocks used is an implementation choice.

For example:

```text
4 blocks
5 blocks
8 blocks
10 blocks
```

can all be used.

Using more blocks can provide a more stable estimate if enough ciphertext is available.

---

# Step 2 — Select Candidate Key Sizes

Instead of trusting only one key size, sort the candidates by normalized Hamming distance.

For example:

```text
Key Size    Normalized Distance
--------------------------------
5           2.81
3           2.94
10          3.02
7           3.15
14          3.20
...
```

The smallest values become candidates for further testing.

This is important because:

> A multiple of the real key size can also produce a good score.

For example, if the actual key size is:

```text
3
```

then:

```text
6
9
12
...
```

may also appear promising.

Therefore, we should test several candidate key sizes rather than blindly accepting the first result.

---

# Step 3 — Transpose the Ciphertext

Once we have a candidate key size, we rearrange the ciphertext.

Suppose:

```text
Key Size = 3
```

and ciphertext positions are:

```text
C0 C1 C2 C3 C4 C5 C6 C7 C8
```

The repeating key is:

```text
K0 K1 K2 K0 K1 K2 K0 K1 K2
```

Therefore:

```text
C0 → K0
C1 → K1
C2 → K2
C3 → K0
C4 → K1
C5 → K2
C6 → K0
C7 → K1
C8 → K2
```

We group bytes according to the key position:

```text
Group 0:

C0 C3 C6


Group 1:

C1 C4 C7


Group 2:

C2 C5 C8
```

This operation is called **transposition**.

---

# Why Transpose?

This is the key insight of the challenge.

Originally, we have:

```text
Repeating-key XOR
```

which looks difficult because multiple key bytes are involved.

After transposition:

```text
Group 0 → encrypted with only K0
Group 1 → encrypted with only K1
Group 2 → encrypted with only K2
```

Each group is now effectively a **single-byte XOR cipher**.

And we already know how to break single-byte XOR from Challenge 3.

---

# Using the Modulo Operator

The important line is:

```go
index := i % keySize
```

For:

```text
keySize = 3
```

we get:

```text
i       i % 3
----------------
0       0
1       1
2       2
3       0
4       1
5       2
6       0
7       1
8       2
```

Therefore:

```text
0 → Group 0
1 → Group 1
2 → Group 2
3 → Group 0
4 → Group 1
5 → Group 2
...
```

This gives us:

```text
Group 0 → C0 C3 C6 C9 ...
Group 1 → C1 C4 C7 C10 ...
Group 2 → C2 C5 C8 C11 ...
```

In other words:

```text
i % keySize
```

tells us **which key byte was used to encrypt the current ciphertext byte**.

---

# Step 4 — Break Each Group

Now each group can be treated as a single-byte XOR ciphertext.

For example:

```text
Group 0 → C0 C3 C6 C9 ...
```

was encrypted using:

```text
K0
```

We try all 256 possible byte values:

```text
0x00
0x01
0x02
...
0xFF
```

For every candidate key byte:

```text
plaintext = ciphertext XOR candidate_key
```

Then score the resulting plaintext for English-like characteristics.

This is exactly the technique used in **Challenge 3**.

---

# Recover the Complete Key

Suppose the three groups produce:

```text
Group 0 → Key byte = K0
Group 1 → Key byte = K1
Group 2 → Key byte = K2
```

We combine them:

```text
Key = K0 K1 K2
```

For example:

```text
K0 = I
K1 = C
K2 = E
```

Therefore:

```text
Key = ICE
```

---

# Step 5 — Decrypt the Ciphertext

Once we have the key, decrypt using the same XOR operation:

```text
Plaintext[i] =
    Ciphertext[i] XOR Key[i % KeySize]
```

XOR has an important property:

```text
A XOR B XOR B = A
```

Therefore, the same operation used for encryption can be used for decryption.

---

# Complete Attack Strategy

The complete process is:

```text
1. Decode the ciphertext
        ↓
2. Try key sizes from 2 to 40
        ↓
3. Divide ciphertext into blocks
        ↓
4. Calculate Hamming distances
        ↓
5. Normalize distances by key size
        ↓
6. Sort key-size candidates
        ↓
7. Take several promising key sizes
        ↓
8. Transpose ciphertext for each candidate
        ↓
9. Break each transposed group as single-byte XOR
        ↓
10. Reconstruct the key
        ↓
11. Decrypt the entire ciphertext
        ↓
12. Score the resulting plaintext
        ↓
13. Select the best candidate
```

---

# Important Concepts

### Hamming Distance

Measures how many **bits differ** between two byte sequences.

```text
XOR → count 1 bits
```

---

### Normalization

Allows Hamming distances from different key sizes to be compared:

```text
distance / keySize
```

---

### Transposition

Rearranges ciphertext bytes so that bytes encrypted with the same key byte are grouped together.

```text
C0 C1 C2 C3 C4 C5
 ↓  ↓  ↓  ↓  ↓  ↓

C0 C3
C1 C4
C2 C5
```

---

### Single-Byte XOR

After transposition, each group is encrypted using only one key byte.

Therefore, Challenge 3 can be reused.

---

### Repeating-Key XOR

Encryption/decryption uses:

```text
data[i] XOR key[i % keySize]
```

---

# Key Insight

The most important idea in this challenge is:

> **A repeating-key XOR cipher can be transformed into multiple independent single-byte XOR problems by grouping ciphertext bytes according to their position within the repeating key.**

So Challenge 6 builds directly on Challenge 3:

```text
Challenge 3
Single-byte XOR
        ↓
Challenge 6
Repeating-key XOR
        ↓
Find key size
        ↓
Transpose ciphertext
        ↓
Multiple single-byte XOR problems
        ↓
Reuse Challenge 3
        ↓
Recover complete key
```

This is the main cryptographic insight behind the challenge.
