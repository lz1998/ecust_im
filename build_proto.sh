#!/usr/bin/env bash

mkdir dto

protoc -I ecust_im_idl --gofast_out=dto ecust_im_idl/*.proto
