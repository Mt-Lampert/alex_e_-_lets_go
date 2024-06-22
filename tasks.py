from invoke import run, task

@task
def test(c):
    print("Invoke is up and running!")

@task
def newSnippet(c):
    cmd = "http --form --ignore-stdin POST 'localhost:3000/create/snippet' title='Currywurst' content='Gehste inne Stadt\\nWat macht dich da satt?\\n’Ne Currywurst!\\n\\n-- Herbert Grönemeyer' expires='1 month'"
    c.run(cmd)
