// Copyright 2020 Matthias Behr
//
// This file is part of libSML.
//
// libSML is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// libSML is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with libSML.  If not, see <http://www.gnu.org/licenses/>.

#include <sml/sml_file.h>

int LLVMFuzzerTestOneInput(const uint8_t *data, size_t size) {
	// we shift by 8 to be able to use libsml-testing bin files inside
	// corpus dir
	sml_file *file = sml_file_parse(size > 8 ? ((unsigned char *)data) + 8 : (unsigned char *)data,
									size > 8 ? (size - 8) : size);
	sml_file_free(file);
	return 0;
}

// compile with:
// clang -g -O1 -L./sml/lib -fsanitize=fuzzer,address,undefined test/src/fuzzer.c `find sml/src/*.c`
//  -luuid -I ./sml/include

// run with:
// mkdir corpus_dir
// prefill corpus_dir e.g. with files libsml-testing
// ./a.out corpus_dir
// run until it stops with a crash file being generated
// e.g. crash-...
// debug with e.g.
// gdb --args ./a.out crash-...
