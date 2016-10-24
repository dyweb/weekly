# HOWTO

```bash
$ cd .jekyll
$ bundle install
$ ./build
```

## Update posts

```bash
$ rm -rf $(seq 2015 `date +%Y`)
$ # cp the posts to root directory
$ ./scripts/add-headers.sh
$ ./.jekyll/build
$ git push
```