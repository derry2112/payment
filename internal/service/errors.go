package service

import "errors"

var ErrProductNotFound = errors.New("produk tidak ditemukan")
var ErrNoProductChanges = errors.New("tidak ada data produk yang diperbarui")
var ErrCategoryNotFound = errors.New("kategori tidak ditemukan")
