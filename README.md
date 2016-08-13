# Justasec

Justasec waits until it is replaced and then executes the replacement with the same arguments.

Have you ever accidentally ran an old build result while a new one was compiling? Add this one line to the beginning of your build script and say goodbye to that mistake!

```shell
cp `which justasec` build/yourprogram
```

That way, every time you start a build, the target will be replaced with a copy of justasec. When your build finishes, it will overwrite the target with your program. The magic is this. You can run your program before it's done building. Justasec will block execution until its file is replaced with your program, and then it will exec your program with the same arguments.
