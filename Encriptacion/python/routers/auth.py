from fastapi import APIRouter, Depends, HTTPException, status
from fastapi.security import OAuth2PasswordRequestForm

router = APIRouter(tags=["auth"])


@router.post("/token")
async def login(form_data: OAuth2PasswordRequestForm = Depends()):
    try:
        uname = (form_data.username or "").strip().lower()
        print("Your username is:", uname)   
        pwd = form_data.password or ""
        print("Your password is:", pwd)
        print(f"/token login attempt user={uname!r} grant_type={form_data.grant_type!r}")
        if uname == "aminespinoza@mail.com" and pwd == "polliTeAmo123":
            return {"access_token": "yourtoken", "token_type": "bearer"}
    except Exception as e:
        print("/token error:", e)
    raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Invalid credentials")
