#!/usr/bin/env perl
use strict;
use warnings;

use Furl;
use JSON::PP;
use Encode;
use Archive::Tar;
use IO::Uncompress::Gunzip qw(:all);

my $client = Furl->new;
my $res = $client->get('https://api.github.com/repos/cli/cli/releases/latest');
unless ($res->is_success) {
    die "Failed to download. " + $res->status_line;
}

my $json = decode_json(decode_utf8 $res->content);
printf "Download version %s\n", $json->{tag_name};

for my $asset (@{$json->{assets}}) {
    my $url = $asset->{browser_download_url};
    if ($url =~ m/linux_amd64\.tar\.gz$/) {
        $res = $client->get($url);
        unless ($res->is_success) {
            die "Failed to $url";
        }

        my $buffer;
        gunzip \$res->content, \$buffer or die "cannot gunzip: $GunzipError";

        open my $fh, '<:raw', \$buffer or die "can't open as string";

        my $tar = Archive::Tar->new;
        $tar->read($fh) or die "cannot read tar file: $!";

        for my $file ($tar->list_files()) {
            if ($file =~ m{/bin/gh$}) {
                $tar->extract_file($file, 'gh') or die "cannot extract 'gh': $!";
                print "Download and extract 'gh'\n";
                exit 0;
            }
        }

        close $fh;
    }
}
