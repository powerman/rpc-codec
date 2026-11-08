# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.3.0] - 2026-11-08

### 🔔 Changed

- Require Go 1.26 by @powerman in [1f28086]

### 🐛 Fixed

- Avoid recursive loop when unmarshalling server request by @peternose in [955adeb]
- Support GOEXPERIMENT=jsonv2 and Go 1.27 by @powerman in [d061d76]

[1.3.0]: https://github.com/powerman/rpc-codec/compare/v1.2.2..v1.3.0
[955adeb]: https://github.com/powerman/rpc-codec/commit/955adeb1013696db824d5ebb21efeea193d21294
[d061d76]: https://github.com/powerman/rpc-codec/commit/d061d7691ed6c691a1421637552e090267e605a1
[1f28086]: https://github.com/powerman/rpc-codec/commit/1f2808669bf7597d622c55f55e7277080d807c64

## [1.2.2] - 2020-04-23

### 🐛 Fixed

- Improve EOF detection by @powerman in [6f3efa7]

[1.2.2]: https://github.com/powerman/rpc-codec/compare/v1.2.1..v1.2.2
[6f3efa7]: https://github.com/powerman/rpc-codec/commit/6f3efa778da423a2df7f44881f5617ae8fff3e84

## [1.2.1] - 2020-04-22

### 🐛 Fixed

- Weird go panic by @powerman in [c2902e8]

[1.2.1]: https://github.com/powerman/rpc-codec/compare/v1.2.0..v1.2.1
[c2902e8]: https://github.com/powerman/rpc-codec/commit/c2902e88d7b3fa7ddd0c3b93ea95210b246a1aab

## [1.2.0] - 2020-04-22

### 🚀 Added

- Add WrapError by @powerman in [e9498a9]

[1.2.0]: https://github.com/powerman/rpc-codec/compare/v1.1.4..v1.2.0
[e9498a9]: https://github.com/powerman/rpc-codec/commit/e9498a983f35578906bb732496a38efa5aa1e4e6

## [1.1.4] - 2020-04-07

### 🐛 Fixed

- Always return non-nil context by @powerman in [6776687]

[1.1.4]: https://github.com/powerman/rpc-codec/compare/v1.1.3..v1.1.4
[6776687]: https://github.com/powerman/rpc-codec/commit/677668791630b3851bb65c5710323607b5b9fcf0

## [1.1.3] - 2019-08-17

[1.1.3]: https://github.com/powerman/rpc-codec/compare/v1.1.2..v1.1.3

## [1.1.2] - 2019-05-24

### 🔔 Changed

- Fix typo in README by @powerman in [4e25cc6]

### 🐛 Fixed

- Jsonrpc2client -version by @powerman in [c72a4f8]

[1.1.2]: https://github.com/powerman/rpc-codec/compare/v1.1.1..v1.1.2
[4e25cc6]: https://github.com/powerman/rpc-codec/commit/4e25cc6723bdb5453d0416c9668e29cc93044569
[c72a4f8]: https://github.com/powerman/rpc-codec/commit/c72a4f8142c5cba1107b0973aedb0ab6f5513d9a

## [1.1.1] - 2019-05-24

### 🔔 Changed

- Update README by @powerman in [5dd47f5]

[1.1.1]: https://github.com/powerman/rpc-codec/compare/v1.1.0..v1.1.1
[5dd47f5]: https://github.com/powerman/rpc-codec/commit/5dd47f59c3ffc4a886c753fa7e28ba592cfd2d13

## [1.1.0] - 2019-05-24

### 🚀 Added

- Add cmd/jsonrpc2client by @powerman in [fe16ac0]

[1.1.0]: https://github.com/powerman/rpc-codec/compare/v1.0.1..v1.1.0
[fe16ac0]: https://github.com/powerman/rpc-codec/commit/fe16ac082cfe744fa3f5a3cc316b0923db8447cf

## [1.0.1] - 2018-11-14

### 🔔 Changed

- Add go.mod by @powerman in [1f88aa2]
- Setup linter by @powerman in [bdee7fe]
- Switch to CircleCI by @powerman in [6d10123]
- Update README by @powerman in [af7f4ca]

[1.0.1]: https://github.com/powerman/rpc-codec/compare/v1.0.0..v1.0.1
[1f88aa2]: https://github.com/powerman/rpc-codec/commit/1f88aa249cd7b42e936f546f56e719200eb00f80
[bdee7fe]: https://github.com/powerman/rpc-codec/commit/bdee7fe61026f205147ec8f777a8b2c973ceed03
[6d10123]: https://github.com/powerman/rpc-codec/commit/6d10123e6e4842b187ff5dfe77a6ece5b37b3e66
[af7f4ca]: https://github.com/powerman/rpc-codec/commit/af7f4ca4cb4a72587df9f428a9226393eff61107

## [1.0.0] - 2018-03-01

### 🔔 Changed

- Start with net/rpc/jsonrpc by @powerman in [27a1ae8]
- Rename package by @powerman in [b8ddb19]
- Add version by @powerman in [ee217b7]
- Add `protocol_test.go` by @powerman in [43b5ead]
- Add Error object support by @powerman in [dba7a1b]
- Add by-position and by-name params by @powerman in [e9cbb9e]
- Add notifications by @powerman in [28fa3d7]
- Improve error when ReadResponseBody fail by @powerman in [a4c8ded]
- Update doc by @powerman in [5b9a312]
- Add batch support for server by @powerman in [54e751d]
- Update copyright by @powerman in [b4de638]
- Move to subdir by @powerman in [7b81936]
- Create LICENSE by @powerman in [206ba3a]
- Add .travis.yml by @powerman in [bf534d5]
- Add README by @powerman in [5f88831]
- Happy lint by @powerman in [bd19b44]
- Add benchmarks by @powerman in [80f9612]
- Implement server-side of "over HTTP" spec by @powerman in [c620de9]
- Use httptest by @powerman in [ed7abc6]
- Implement client-side of "over HTTP" spec, close #1 by @powerman in [fd660d2]
- Update README by @powerman in [737d776]
- Rename DialHTTP to NewHTTPClient by @powerman in [ad9dc3a]
- Improve doc by @powerman in [eed1bac]
- HTTP Client: allow parallel requests by @powerman in [24a2f6d]
- Add coveralls by @powerman in [afe254f]
- Work around golang/go#12513 by @powerman in [37e42af]
- Work around golang/go#12262 by @powerman in [65bac2c]
- Cleanup .travis.yml by @powerman in [7caf9ad]
- Return errParse on incorrect HTTP, fix #4 by @powerman in [d04fa12]
- Simplify tests by @powerman in [d3b4477]
- Improve doc, add example; close #2 by @powerman in [dd1bd07]
- Update doc by @powerman in [bccf0e6]
- Add custom HTTP client support by @powerman in [f64d347]
- Enable HTTP conn reuse after errors by @powerman in [2f8d743]
- Add context support by @powerman in [ffdd59f]
- Require Go 1.7 because of context by @powerman in [7815806]
- Document context support by @powerman in [f78561a]
- Add context support for HTTPHandler by @powerman in [7232415]
- Add tests for context (except batch) by @powerman in [817f47b]
- Add test for just fixed bug by @powerman in [dd19550]
- Add test for batch with context by @powerman in [e95ddc1]
- Use go 1.8 by @powerman in [b190271]
- Add test to reproduce charset issue by @anpryl in [9bbcecf]
- Change header parsing to fix charset issue by @anpryl in [ae2efb2]
- Fix contentType in test by @anpryl in [b7a13a6]
- Lowercase standard error messages by @powerman in [91cc933]
- Handle ParseMediaType error by @powerman in [761c1da]
- Improve test by @powerman in [89c9f32]
- Exported NewClientCodec. by @trey-jones in [5c155bf]
- Allow use of a client codec other than the one supplied this package for client. by @trey-jones in [72c2c29]
- Updated original new client method to use NewClientWithCodec by @trey-jones in [c217c72]

### 🐛 Fixed

- Fix protocol validation by @powerman in [0e6169e]
- Fix tests for 1.5 by @powerman in [9da8bd8]
- Fix build for go1.1-1.2 by @powerman in [0077900]
- Fix build for go1.1-1.2 by @powerman in [e98f696]
- Fix crash on go-1.1 and go-1.5 by @powerman in [8fcc69a]
- Fix .travis.yml by @powerman in [d0a6e70]
- Fix tests for go-1.1 by @powerman in [e40d686]
- Fix example by @powerman in [103c760]
- Fix tests for Go tip by @powerman in [b12587c]
- Fix absent context for method without params by @powerman in [81335e6]
- Fix response content type checking by Alexey Makhov in [ba76484]
- Fix response content type checking by Alexey Makhov in [c38fe1b]
- Fix response content type checking. tests by Alexey Makhov in [8b64234]
- Fix typo by @powerman in [c92d5b2]
- Fix test by @powerman in [ea81e3d]
- Fix goroutine leak on HTTP Client.Close by @powerman in [7e75941]

[1.0.0]: https://github.com/powerman/rpc-codec/compare/%40%7B10year%7D..v1.0.0
[27a1ae8]: https://github.com/powerman/rpc-codec/commit/27a1ae8e983040df4c37e401e01312333b479b17
[b8ddb19]: https://github.com/powerman/rpc-codec/commit/b8ddb19312fe212b0c7c4df8ef2efe43ad99af1a
[ee217b7]: https://github.com/powerman/rpc-codec/commit/ee217b7a5c111eefb89a6460609ca981e1de13e7
[43b5ead]: https://github.com/powerman/rpc-codec/commit/43b5eadefe2a910472f52bf5f671bf6d67c4896c
[dba7a1b]: https://github.com/powerman/rpc-codec/commit/dba7a1bb3dc5ede2f500cd02ae8ef5513c5a1457
[e9cbb9e]: https://github.com/powerman/rpc-codec/commit/e9cbb9e33abeebce3c55aed37ac49c958014a1cc
[28fa3d7]: https://github.com/powerman/rpc-codec/commit/28fa3d7f1078817c40dc47e2933a3ce21f31c221
[0e6169e]: https://github.com/powerman/rpc-codec/commit/0e6169eb45e98f0f8487e51eb07220cb50649957
[a4c8ded]: https://github.com/powerman/rpc-codec/commit/a4c8ded75db29eb2dc58676d643fe54759cfdfe6
[5b9a312]: https://github.com/powerman/rpc-codec/commit/5b9a312d9118a1051d4552f2ef855a3016629f54
[54e751d]: https://github.com/powerman/rpc-codec/commit/54e751daddfcc39c451167be4afec014f0786e8f
[b4de638]: https://github.com/powerman/rpc-codec/commit/b4de638814f21afe75cd629531075cb08759d4d2
[7b81936]: https://github.com/powerman/rpc-codec/commit/7b819366b5d0663a01472f603ec3180360e47dd3
[9da8bd8]: https://github.com/powerman/rpc-codec/commit/9da8bd870668101ef2f8074dedd4ca7b5d60b576
[206ba3a]: https://github.com/powerman/rpc-codec/commit/206ba3af32654deef0c818c1629d2458696506e1
[bf534d5]: https://github.com/powerman/rpc-codec/commit/bf534d5b32bc3cc80e7d4d837b4d14c6e90a77fa
[5f88831]: https://github.com/powerman/rpc-codec/commit/5f88831a60fc690c537fe44434707ba881a7a41d
[bd19b44]: https://github.com/powerman/rpc-codec/commit/bd19b44802de819d5f25448706e70a64d1f4f1a2
[80f9612]: https://github.com/powerman/rpc-codec/commit/80f96126329a1508dd78b359a50479d869f1fcd0
[c620de9]: https://github.com/powerman/rpc-codec/commit/c620de92da0f24e8696cc5fce412b670a3ac856b
[0077900]: https://github.com/powerman/rpc-codec/commit/00779002a53b2d7f96e386b958ef5117ac906653
[ed7abc6]: https://github.com/powerman/rpc-codec/commit/ed7abc66372f3887a56be9046b30dbda1328df1e
[fd660d2]: https://github.com/powerman/rpc-codec/commit/fd660d2f4713790152a0c625417999343c6bdfc1
[737d776]: https://github.com/powerman/rpc-codec/commit/737d776c03dc535e9c3c7dfdcdd5e0c60f8c27f7
[e98f696]: https://github.com/powerman/rpc-codec/commit/e98f6960b8bd9a99422d193064ec3930956eecb8
[ad9dc3a]: https://github.com/powerman/rpc-codec/commit/ad9dc3ae23c49a0783ded3a168d32fdd973513c0
[eed1bac]: https://github.com/powerman/rpc-codec/commit/eed1bac809c82060cb0d61d7c367ee26de5354e8
[24a2f6d]: https://github.com/powerman/rpc-codec/commit/24a2f6d991b6e748aa5bccd45808a50b2bf19fc0
[afe254f]: https://github.com/powerman/rpc-codec/commit/afe254f86a74dd91d8a48a1f89319c040851052a
[37e42af]: https://github.com/powerman/rpc-codec/commit/37e42af56037f6a0b6715b1f5fb86d98b610fe4e
[8fcc69a]: https://github.com/powerman/rpc-codec/commit/8fcc69a74efc681f19e1f1bedff98e06306794f3
[d0a6e70]: https://github.com/powerman/rpc-codec/commit/d0a6e709d4dc1fd0967805bca3ca7f54b921b0ad
[65bac2c]: https://github.com/powerman/rpc-codec/commit/65bac2c4aa55298657dfb199f1452afe10d8183b
[7caf9ad]: https://github.com/powerman/rpc-codec/commit/7caf9adb031d6fba8081b6914e6d9739356979e1
[d04fa12]: https://github.com/powerman/rpc-codec/commit/d04fa1245adef5e98074978c9487e864e2290060
[e40d686]: https://github.com/powerman/rpc-codec/commit/e40d68678e471d6be97b8083c68d7c2df516bb32
[d3b4477]: https://github.com/powerman/rpc-codec/commit/d3b447793c267710a52dee05b0fd6d0148012c32
[dd1bd07]: https://github.com/powerman/rpc-codec/commit/dd1bd07b01622e6f876a1b33d21a70091210939a
[103c760]: https://github.com/powerman/rpc-codec/commit/103c7609a4eacb1f8221c73cfecd3d4e9003e21b
[bccf0e6]: https://github.com/powerman/rpc-codec/commit/bccf0e672639ac66adea03fe8c294a8c9490e3ab
[f64d347]: https://github.com/powerman/rpc-codec/commit/f64d34739e8d906b04ca02432a3f736a0025a5f9
[2f8d743]: https://github.com/powerman/rpc-codec/commit/2f8d7434e99ddbac66da0f4b5d68c1591cf9eab1
[b12587c]: https://github.com/powerman/rpc-codec/commit/b12587c9f0f1616a11abaa9f576238176fef2740
[ffdd59f]: https://github.com/powerman/rpc-codec/commit/ffdd59fd5e21b0eebdabe94272d263ddb7e7f01e
[7815806]: https://github.com/powerman/rpc-codec/commit/781580618225c584a5ded7e6e7ebffd3436830a1
[f78561a]: https://github.com/powerman/rpc-codec/commit/f78561aafaac5c5141640cbd00ce11ba9b9d1395
[7232415]: https://github.com/powerman/rpc-codec/commit/7232415c8be827d300049efa5bd7af8fbaad393a
[817f47b]: https://github.com/powerman/rpc-codec/commit/817f47b7df16630f95f76b93a00191f2ba44b7bc
[81335e6]: https://github.com/powerman/rpc-codec/commit/81335e6ee237a7532a31eef13c32f4fcb261c37b
[dd19550]: https://github.com/powerman/rpc-codec/commit/dd19550c8d75c58af8cee2764ab8eec79fcf46d8
[e95ddc1]: https://github.com/powerman/rpc-codec/commit/e95ddc17573826f0428a136e5b06211e23013f93
[b190271]: https://github.com/powerman/rpc-codec/commit/b190271b7894e3db2143ba10c09705d1d19f8ffe
[9bbcecf]: https://github.com/powerman/rpc-codec/commit/9bbcecfb47ffa1963dac388888b1ddc5ce6db71c
[ae2efb2]: https://github.com/powerman/rpc-codec/commit/ae2efb2d7963ca1cb78f9f6bfd6978b8854791cf
[b7a13a6]: https://github.com/powerman/rpc-codec/commit/b7a13a6e46aac14a47eeece1e38d37e6926263fb
[91cc933]: https://github.com/powerman/rpc-codec/commit/91cc9337d50296e4b44af14c84e9ac4739fe015a
[ba76484]: https://github.com/powerman/rpc-codec/commit/ba76484a06c309faf0a2849be32e953d159240ea
[c38fe1b]: https://github.com/powerman/rpc-codec/commit/c38fe1b8fe872656f1f5af44395a8c2d0616dadd
[8b64234]: https://github.com/powerman/rpc-codec/commit/8b64234acd19e1458ada45c04c20e04e8e9c01ee
[c92d5b2]: https://github.com/powerman/rpc-codec/commit/c92d5b2684eb39be05798e9463e3d52dacbdc1ed
[761c1da]: https://github.com/powerman/rpc-codec/commit/761c1daa66f8b5c57b28da0cf02a22ecbca5a895
[89c9f32]: https://github.com/powerman/rpc-codec/commit/89c9f324b0410420d65936124f4cf78dbc360e6e
[ea81e3d]: https://github.com/powerman/rpc-codec/commit/ea81e3d929d602f114634c2e7ddc510104d291af
[5c155bf]: https://github.com/powerman/rpc-codec/commit/5c155bf10179a737455b0ce44696850f9d15dd3e
[72c2c29]: https://github.com/powerman/rpc-codec/commit/72c2c29752f0c9f9f1e7e3a04f0ee7dd508fd8cf
[c217c72]: https://github.com/powerman/rpc-codec/commit/c217c720eff1b83e28949ab895e2ad39da4810bb
[7e75941]: https://github.com/powerman/rpc-codec/commit/7e759419863c5bce9aaad2c1bd3236a3a9dc35e4

<!-- generated by git-cliff -->
