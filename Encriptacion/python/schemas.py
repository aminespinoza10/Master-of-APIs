from pydantic import BaseModel

class User(BaseModel):
    id: int
    name: str
    username: str
    password: str


class CreateUser(BaseModel):
    name: str
    username: str
    password: str


class UserPublic(BaseModel):
    id: int
    name: str
    username: str