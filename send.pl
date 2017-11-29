package send;

use strict;
use warnings;

use Data::Dumper;
use Try::Tiny;
use Mojo::UserAgent;
use JSON qw(decode_json);
use Scalar::Util qw(looks_like_number);
use IO::Async::Loop;
use IO::Async::Timer::Periodic;

# send request every 3 seconds
my $loop = IO::Async::Loop->new;
my $timer = IO::Async::Timer::Periodic->new(
    interval => 3,
    on_tick => sub { get_sum(rand(1000),rand(1000)/$_) for 1..1000 },
);
$timer->start;
$loop->add($timer);
$loop->run;

sub get_sum {
    my ($a, $b) = @_;

    my $tx = Mojo::UserAgent->new->get("::3000/sum?a=$a&b=$b");
    my $resp = $tx->result->content->asset;
    my $result = decode_json($resp->{content})->{Result};

    return $result + 0;
}

1;
