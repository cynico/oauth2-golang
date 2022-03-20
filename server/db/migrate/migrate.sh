#!/bin/sh

migrate -database mysql://root:dbpass@/oauth -path db/migrate/mysql up