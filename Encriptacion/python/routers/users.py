from fastapi import APIRouter, Depends, HTTPException
from fastapi.responses import PlainTextResponse
from fastapi.security import OAuth2PasswordBearer
from typing import List
import os
import asyncpg
import bcrypt

from schemas import CreateUser, UserPublic

router = APIRouter(tags=["users"])
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")


@router.get("/okCode", response_class=PlainTextResponse, name="getOkCode")
async def get_ok_code(_: str = Depends(oauth2_scheme)):
	return PlainTextResponse("Everything is awesome!", status_code=200)


@router.get("/getUsers", response_model=List[UserPublic])
async def get_users(_: str = Depends(oauth2_scheme)):
	db_url = os.getenv("DATABASE_URL")
	if not db_url:
		raise HTTPException(status_code=500, detail="DATABASE_URL not set")

	try:
		conn = await asyncpg.connect(db_url)
		rows = await conn.fetch("SELECT id, name, username FROM users")
	except Exception:
		raise HTTPException(status_code=500, detail="DB error")
	finally:
		if 'conn' in locals():
			await conn.close()

	return [UserPublic(id=row["id"], name=row["name"], username=row["username"]) for row in rows]


@router.post("/createUser", response_model=UserPublic)
async def create_user(payload: CreateUser, _: str = Depends(oauth2_scheme)):
	db_url = os.getenv("DATABASE_URL")
	if not db_url:
		raise HTTPException(status_code=500, detail="DATABASE_URL not set")

	# Hash password with bcrypt before storing
	if not payload.password or not payload.password.strip():
		raise HTTPException(status_code=400, detail="Password required")
	hashed_pw = bcrypt.hashpw(payload.password.encode("utf-8"), bcrypt.gensalt()).decode("utf-8")

	try:
		conn = await asyncpg.connect(db_url)
		row = await conn.fetchrow(
			"INSERT INTO users (name, username, password) VALUES ($1, $2, $3) RETURNING id, name, username",
			payload.name, payload.username, hashed_pw,
		)
	except Exception:
		raise HTTPException(status_code=500, detail="DB error")
	finally:
		if 'conn' in locals():
			await conn.close()

	return UserPublic(id=row["id"], name=row["name"], username=row["username"])

