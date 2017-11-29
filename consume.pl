#!/usr/bin/env perl
package consume;

use strict;
use warnings;

use Try::Tiny;
use Data::Dumper;
use JSON qw(decode_json);
use AnyEvent::RabbitMQ;

# Get notified on events in sum queue
my $cv = AnyEvent->condvar;
my $conn = AnyEvent::RabbitMQ->new->load_xml_spec()->connect(
    vhost => '/',
    host => '127.0.0.1',
    port => 5672,
    user => 'guest',
    pass => 'guest',
    on_success => sub {
        shift->open_channel(
            on_success => sub {
                my $ch = shift;
                print("[*] Succesfully connected to AMQP server\n");
                $ch->declare_queue(
                    queue => 'sum',
                    on_success => sub {
                        print("[*] Queue `sum` declared, start listening\n");
                        $ch->consume(on_consume => \&parse_message, no_ack => 1);
                    },
                    on_failure => $cv,
                );
            },
            on_failure => $cv,
            on_close => \&fatal_conn_close,
        );
    },
    on_failure => $cv,
    on_read_failure => sub { die @_ },
    on_close => \&fatal_conn_close,
);

# Log RabbitMQ close connection reason
sub fatal_conn_close {
    my $frame = shift->method_frame;
    die "[x] Connection closed", $frame->reply_code, $frame->reply_text;
}

sub parse_message {
    my $result = decode_json(shift->{body}->{payload})->{Result};
    # make some load
    my @brainfuck = map { sqrt($result) * rand($_) } 0..rand(100);
    print "[x] Received: " . Dumper(\@brainfuck) . "\n";
}

$cv->recv;

1;
