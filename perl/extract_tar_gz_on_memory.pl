#!/usr/bin/env perl
use strict;
use warnings;

use HTTP::Tiny;
use Encode;
use Archive::Tar;
use IO::Uncompress::Gunzip qw(:all);

my $dest = shift;
my $url = 'https://www.cpan.org/src/5.0/perl-5.36.0.tar.gz';

my $res = HTTP::Tiny->new->get($url);
unless ($res->{success}) {
    die "$res->{status} $res->{reason}";
}

my $tar_buf;
gunzip \$res->{content}, \$tar_buf or die "gunzip failed: $GunzipError\n"; 

open my $fh, '<', \$tar_buf or die "Can't open tar content as a file handle: $!";

my $tar = Archive::Tar->new;
$tar->read($fh);

if (defined $dest) {
    $tar->setcwd($dest);
}

$tar->extract;

close $fh;